package models

import ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"

type ForumAndThreads struct {
	// Forum parameters.
	ForumId        uint   `json:"forumId"`
	ForumSectionId uint   `json:"forumSectionId"`
	ForumName      string `json:"forumName"`
	EventData

	// Thread parameters. If pagination is used, these lists contain
	// information after the application of pagination.
	ThreadIds ul.UidList `json:"threadIds"`
	Threads   []Thread   `json:"threads"`

	// Number of the current page of threads.
	Page *uint `json:"page,omitempty"`

	// Total number of available pages of threads.
	TotalPages *uint `json:"totalPages,omitempty"`

	// Total number of available threads.
	TotalThreads *uint `json:"totalThreads,omitempty"`
}

func NewForumAndThreads(forum *Forum) (fat *ForumAndThreads) {
	if forum == nil {
		return nil
	}

	fat = &ForumAndThreads{
		ForumId:        forum.Id,
		ForumSectionId: forum.SectionId,
		ForumName:      forum.Name,

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

	return fat
}
