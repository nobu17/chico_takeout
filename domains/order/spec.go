package order

type FoodItemRemainQuantitySpecification struct {
	qMap *quantityMap
}

func newFoodItemRemainQuantitySpecification(sameDateOrders []OrderInfo) *FoodItemRemainQuantitySpecification {
	spec := FoodItemRemainQuantitySpecification{}
	spec.calcQuantityMap(sameDateOrders)
	return &spec
}

func (f *FoodItemRemainQuantitySpecification) calcQuantityMap(sameDateOrders []OrderInfo) {
	f.qMap = newQuantityMap()
	for _, order := range sameDateOrders {
		for _, food := range order.GetFoodItems() {
			f.qMap.Add(food.GetItemId(), food.GetQuantity())
		}
	}
}

func (s *FoodItemRemainQuantitySpecification) IsOverRemain(id string, quantity, maxOrderPerDay int) bool {
	current := s.qMap.GetQuantity(id)
	return current+quantity > maxOrderPerDay
}

type quantityMap struct {
	maps map[string]int
}

func newQuantityMap() *quantityMap {
	return &quantityMap{
		maps: map[string]int{},
	}
}

func (q *quantityMap) Add(id string, quantity int) {
	_, ok := q.maps[id]
	if !ok {
		q.maps[id] = quantity
		return
	}
	q.maps[id] = q.maps[id] + quantity
}
func (q *quantityMap) GetQuantity(id string) int {
	_, ok := q.maps[id]
	if !ok {
		return 0
	}
	return q.maps[id]
}
