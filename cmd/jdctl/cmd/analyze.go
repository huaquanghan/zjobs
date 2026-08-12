package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"zjobs/internal/analysis"
	"zjobs/internal/config"
	"zjobs/internal/domain"
	"zjobs/internal/ingest"
	"zjobs/internal/reporting"
)

// jdFlags are the locked input surface for `analyze jd`.
// Exactly one of --url, --file, or --paste must be provided; --profile and
// --cv are required paths to the user's profile YAML and CV variant Markdown.
// --analysis points at the JSON payload emitted by the Claude skill (the
// semantic layer); Go validates it and computes the deterministic verdict.
type jdFlags struct {
	profile  string
	cv       string
	url      string
	file     string
	paste    string
	analysis string
	out      string
}

func newAnalyzeCmd() *cobra.Command {
	f := &jdFlags{}

	var analyzeCmd = &cobra.Command{
		Use:   "analyze",
		Short: "Analyze a JD against a profile + CV variant",
		Args:  cobra.NoArgs,
	}

	var jdCmd = &cobra.Command{
		Use:   "jd",
		Short: "Analyze one JD against one profile + CV variant",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.validate(); err != nil {
				return err
			}
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}

			profile, err := domain.LoadProfile(f.profile)
			if err != nil {
				return err
			}
			cv, err := domain.LoadCVVariant(f.cv)
			if err != nil {
				return err
			}

			var jd *domain.JobDescription
			switch {
			case f.url != "":
				jd, err = ingest.FromURL(f.url, nil)
			case f.file != "":
				jd, err = ingest.FromFile(f.file)
			default:
				jd = ingest.FromPaste(f.paste)
			}
			if err != nil {
				return err
			}

			analysisPayload, err := os.ReadFile(f.analysis)
			if err != nil {
				return fmt.Errorf("--analysis: %w", err)
			}
			ana, err := analysis.Validate(analysisPayload)
			if err != nil {
				return err
			}
			res, err := analysis.Evaluate(ana, analysis.Weights{
				analysis.AxisHardConstraints:  cfg.Rubric.HardConstraints,
				analysis.AxisMustHaveSkills:   cfg.Rubric.MustHaveSkills,
				analysis.AxisNiceToHaveSkills: cfg.Rubric.NiceToHaveSkills,
				analysis.AxisSeniorityScope:   cfg.Rubric.SeniorityScope,
				analysis.AxisDomainContext:    cfg.Rubric.DomainContext,
				analysis.AxisEvidenceStrength: cfg.Rubric.EvidenceStrength,
			})
			if err != nil {
				return err
			}

			report := reporting.Build(jd, profile.Name, cv.Name, ana, res)
			jsonPath, mdPath, err := reporting.WriteReport(f.out, report)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "verdict=%s score=%.2f\njson=%s\nmd=%s\n",
				res.Verdict, res.Score, jsonPath, mdPath)
			return nil
		},
	}

	jdCmd.Flags().StringVar(&f.profile, "profile", "", "path to profile YAML (required)")
	jdCmd.Flags().StringVar(&f.cv, "cv", "", "path to CV variant Markdown (required)")
	jdCmd.Flags().StringVar(&f.url, "url", "", "public JD URL (exactly one of --url/--file/--paste)")
	jdCmd.Flags().StringVar(&f.file, "file", "", "path to JD file (exactly one of --url/--file/--paste)")
	jdCmd.Flags().StringVar(&f.paste, "paste", "", "pasted JD text (exactly one of --url/--file/--paste)")
	jdCmd.Flags().StringVar(&f.analysis, "analysis", "", "path to Claude analysis JSON (required)")
	jdCmd.Flags().StringVar(&f.out, "out", "./data/reports", "report output directory")

	analyzeCmd.AddCommand(jdCmd)
	return analyzeCmd
}

func (f *jdFlags) validate() error {
	if f.profile == "" {
		return fmt.Errorf("--profile is required")
	}
	if f.cv == "" {
		return fmt.Errorf("--cv is required")
	}
	if f.analysis == "" {
		return fmt.Errorf("--analysis is required")
	}
	sources := 0
	for _, v := range []string{f.url, f.file, f.paste} {
		if v != "" {
			sources++
		}
	}
	if sources == 0 {
		return fmt.Errorf("exactly one of --url, --file, or --paste is required")
	}
	if sources > 1 {
		return fmt.Errorf("exactly one of --url, --file, or --paste is required; got %d", sources)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(newAnalyzeCmd())
}
