package models

import ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"

type ForumWithThreads struct {
	// Forum parameters.
	ForumId       uint        `json:"forumId"`
	ForumParent   *uint       `json:"forumParent"`
	ForumChildren *ul.UidList `json:"forumChildren"`
	ForumName     string      `json:"forumName"`
	EventData

	// Number of the current page of threads.
	Page *uint `json:"page,omitempty"`

	// Total number of available pages of threads.
	TotalPages *uint `json:"totalPages,omitempty"`

	// Total number of available threads.
	TotalThreads *uint `json:"totalThreads,omitempty"`

	// Thread parameters. If pagination is used, these lists contain
	// information after the application of pagination.
	ThreadIds ul.UidList `json:"threadIds"`
	Threads   []Thread   `json:"threads"`
}

func NewForumWithThreads(forum *Forum) (fwt *ForumWithThreads) {
	if forum == nil {
		return nil
	}

	fwt = &ForumWithThreads{
		ForumId:       forum.Id,
		ForumParent:   forum.Parent,
		ForumChildren: forum.Children,
		ForumName:     forum.Name,

		EventData: EventData{
			Creator: &EventParameters{
				UserId: forum.Creator.UserId,
				Time:   forum.Creator.Time,
			},
			Editor: &OptionalEventParameters{
				UserId: forum.Editor.UserId,
				Time:   forum.Editor.Time,
			},
		},
	}

	return fwt
}
