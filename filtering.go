package word_index

import "log"

type Filtering interface {
	Filter(search *Search) []ItemID
}

type AndOperator struct {
	Left, Right Filtering
}

func NewAndOperator(left, right Filtering) *AndOperator {
	return &AndOperator{Left: left, Right: right}
}

type OrOperator struct {
	Left, Right Filtering
}

func NewOrOperator(left, right Filtering) *OrOperator {
	return &OrOperator{Left: left, Right: right}
}

func (o *AndOperator) Filter(search *Search) []ItemID {
	return mergeOrderedArrayAnd(
		o.Left.Filter(search),
		o.Right.Filter(search),
	)
}

func (o *OrOperator) Filter(search *Search) []ItemID {
	log.Println(o.Left.Filter(search))
	log.Println(o.Right.Filter(search))
	return mergeOrderedArrayOr(
		o.Left.Filter(search),
		o.Right.Filter(search),
	)
}

type FilterIn []Feature

func NewFilterIn(feature ...Feature) FilterIn {
	return feature
}

func (f FilterIn) Filter(search *Search) []ItemID {
	return search.Find(f...)
}
