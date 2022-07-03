export type PerDayOrderableInfo = {
  date: string;
  hourTypeId: string;
  startTime: string;
  endTime: string;
  categories: CategoryItemInfo[];
};

export type CategoryItemInfo = {
  kindId: string;
  title: string;
  priority: number;
  items: OrderableItemInfo[];
};

export type OrderableItemInfo = {
  id: string;
  name: string;
  type: "food" | "stock";
  imageUrl: string;
  memo: string;
  price: number;
  max: number;
};
