package slambox

import "mask_of_the_tomb/internal/core/maths"

// Attachments cannot be collided with, and do not interact
// with anything.
type Attachment struct {
	rect *maths.Rect
}
