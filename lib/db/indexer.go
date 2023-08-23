package db

import (
  "sync"
  "strings"
  "odin/lib/utils"
)

type Document struct {
  ID int
  Title string
  Description string
  Content interface{}
  CreatedAt string
  UpdatedAt string
}

type Index struct {
  docs map[int]*Document
  index map[string][]int
  lock sync.RWMutex
}

func NewIndex()*Index{
  return &Index{
    docs: make(map[int]*Document),
    index: make(map[string][]int),
  }
}
/*
func (idx *Index) AddDocument(doc *Document){
  idx.lock.RLock()
  defer idx.lock.RUnlock()
  idx.docs[doc.ID] = doc
  tokens := strings.Fields(doc.Description)
  for _,token := range tokens {
    if _,ok := idx.index[token]; !ok {
      idx.index[token] = make(idx.index[token],doc.ID,doc.ID)
    }
  }
}
*/

func (idx *Index) Search(query string) []*Document{
  idx.lock.RLock()
  defer idx.lock.RUnlock()
  tokens := strings.Fields(query)
  docIDs := make([][]int,0)
  for _,token := range tokens{
    if ids,ok := idx.index[token];ok {
      docIDs = append(docIDs,ids)
    }
  }
  intersect := Intersect(docIDs)
  documents := make([]*Document,0)
  for _,id := range intersect {
    if doc,ok := idx.docs[id]; ok {
      documents = append(documents,doc)
    }
  }
  return documents
}

var Intersect = func(arr [][]int) []int{
  if len(arr) == 0 {
    return make([]int,0)
  }
  result := make([]int,len(arr[0]))
  copy(result,arr[0])
  for _,array := range arr[1:]{
    newResult := make([]int,0)
    for _,id := range array{
      if utils.ArrayContainsInt(result,id){
        newResult = append(newResult,id)
      }
    }
    result = newResult
  }
  return result
}
