package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	//изначально в структуре было только поле List, чтобы код скомпилировался
	//так как если мы делаем такое встраивание или embedding -> означает, что эта структура реализует интерфейс List,
	//даже если его методы пока не реализованы

	//так как список двусвязный, иным способом до данных не добраться
	first *ListItem
	last  *ListItem
	size  int
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v

	if l.size == 0 {
		l.first = newItem
		l.last = newItem
	} else {
		newItem.Next = l.first
		l.first.Prev = newItem
		l.first = newItem
	}
	l.size++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v

	if l.size == 0 {
		l.first = newItem
		l.last = newItem
	} else {
		newItem.Prev = l.last
		l.last.Next = newItem
		l.last = newItem
	}
	l.size++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i == l.first {
		l.first = i.Next
	}
	if i == l.last {
		l.last = i.Prev
	}

	//"замыкаем" соседей убранного элемента
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	//первый элемент не проверяем, так как он уже в начале списка
	if i == l.first {
		return
	}

	if i == l.last {
		l.last = i.Prev
	}

	//"замыкаем" соседей убранного элемента
	i.Prev.Next = i.Next

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Prev = nil
	i.Next = l.first
	l.first.Prev = i
	l.first = i
}

func NewList() List {
	return new(list) //тип указатель - *list
}
