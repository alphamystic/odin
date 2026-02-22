package definers

/*
  * Defines the blog page data
*/

import (
  "github.com/alphamystic/odin/lib/utils"
)

type Blog struct {
    ID          int       `json:"id"`
    UUID        string    `json:"uuid"`
    OwnerID     string    `json:"owner_id"`
    Title       string    `json:"title"`
    Author      string    `json:"author"`
    Content     string    `json:"content"` // HTML from editor
    MainTag     string    `json:"maintag"`

    // Dynamic JSON arrays (scalable, no schema changes ever)
    Types       []string  `json:"types"`       // e.g ["framework", "report", "policy"]
    Categories  []int     `json:"categories"`  // link to categories table
    SubCats     []int     `json:"subcats"`     // link to categories table
    Tags        []string  `json:"tags"`        // freeform string tags

    Public      bool      `json:"public"`
    Archived    bool      `json:"archived"`
    utils.TimeStamps
}


type Comment struct {
    ID         int    `json:"id"`
    CommentUUID string `json:"comment_uuid"`
    BlogUUID   string `json:"blog_uuid"`
    Comment    string `json:"comment"`
    Commentor  string `json:"commentor"`
    utils.TimeStamps
}

type Category struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    ParentID  *int   `json:"parent_id"` // null = main category, int = subcategory, let null be 1
}


type EnumType struct {
    ID     int    `json:"id"`
    Type   string `json:"type"`
    Description string `json:"description"`
}



// Types Examples
// [
//   "framework",
//   "policy",
//   "evidence",
//   "SOC-report",
//   "IR-SOP"
// ]
