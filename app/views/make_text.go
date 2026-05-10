package views

import (
	"fmt"
	"io"

	"github.com/rajahafify/spritey/app/models"
)

func WriteMakeText(writer io.Writer, result models.MakeResult) error {
	if _, err := fmt.Fprintln(writer, "ok: make"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "png: %s\n", result.Outputs.PNG.Path); err != nil {
		return err
	}
	if result.Outputs.Report != nil {
		if _, err := fmt.Fprintf(writer, "report: %s\n", result.Outputs.Report.Path); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(writer, "frame_count: %d\n", result.Summary.FrameCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "canvas: %dx%d\n", result.Summary.Canvas.Width, result.Summary.Canvas.Height); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "animation_count: %d\n", result.Summary.AnimationCount); err != nil {
		return err
	}
	return nil
}
