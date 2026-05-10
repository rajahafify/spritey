package views

import (
	"fmt"
	"io"

	"github.com/rajahafify/spritey/app/models"
)

func WriteMakeBatchText(writer io.Writer, result models.MakeBatchResult) error {
	if _, err := fmt.Fprintln(writer, "ok: make-batch"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "jobs: %d\n", result.Summary.JobCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "success: %d\n", result.Summary.SuccessCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "failed: %d\n", result.Summary.FailedCount); err != nil {
		return err
	}
	for i, job := range result.Jobs {
		if len(job.Errors) == 0 {
			if _, err := fmt.Fprintf(writer, "job %d: %s ok\n", i+1, job.ID); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprintf(writer, "job %d: %s failed %s\n", i+1, job.ID, job.Errors[0].Code); err != nil {
			return err
		}
	}
	return nil
}
