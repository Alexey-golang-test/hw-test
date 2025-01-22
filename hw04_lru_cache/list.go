package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(vv interface{}) *ListItem
	PushBack(vv interface{}) *ListItem
	Remove(ii *ListItem)
	MoveToFront(ii *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstItem *ListItem
	lastItem  *ListItem
	count     int
}

func NewList() List {
	return new(list)
}

func (list *list) Len() int {
	return list.count
}

func (list *list) Front() *ListItem {
	return list.firstItem
}

func (list *list) Back() *ListItem {
	return list.lastItem
}

func (list *list) PushFront(vv interface{}) *ListItem {
	newItem := &ListItem{vv, list.firstItem, nil}

	// Если уже есть элементы в списке, то у текущего первого элемента надо проставить ссылку на добавляемый элемент
	if list.firstItem != nil {
		list.firstItem.Prev = newItem
	}

	// Если ссылка на последний элемент пустая, то заполним ее первым добавляемым элементом
	if list.lastItem == nil {
		list.lastItem = newItem
	}

	list.firstItem = newItem
	list.count++

	return newItem
}

func (list *list) PushBack(vv interface{}) *ListItem {
	newItem := &ListItem{vv, nil, list.lastItem}

	// Если уже есть элементы в списке, то у текущего последнего элемента надо проставить ссылку на добавляемый элемент
	if list.lastItem != nil {
		list.lastItem.Next = newItem
	}

	// Если ссылка на первый элемент пустая, то заполним ее первым добавляемым элементом
	if list.firstItem == nil {
		list.firstItem = newItem
	}

	list.lastItem = newItem
	list.count++

	return newItem
}

func (list *list) Remove(ii *ListItem) {
	// Есть следующий элемент
	if ii.Next != nil {
		ii.Next.Prev = ii.Prev
	}

	// Есть предыдущий элемент
	if ii.Prev != nil {
		ii.Prev.Next = ii.Next
	}

	// Удаляется первый элемент
	if list.firstItem == ii {
		list.firstItem = ii.Next
	}

	// Удаляется последний элемент
	if list.lastItem == ii {
		list.lastItem = ii.Prev
	}

	list.count--
}

func (list *list) MoveToFront(ii *ListItem) {
	list.Remove(ii)
	list.PushFront(ii.Value)
}
