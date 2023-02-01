import { Cart, OptionItemInfo } from "../../hooks/UseItemCart";

type ItemRequest = {
  item: RequestItemInfo;
  quantity: number;
  selectOptions: OptionItemInfo[];
};

type RequestItemInfo = {
  price: number;
};

type OrderItem = {
  price: number;
  quantity: number;
  options: OrderOptionItem[];
};

type OrderOptionItem = {
  price: number;
};

export const getTotalPrice = (items: OrderItem[]): string => {
  let total = 0;
  items.forEach((item) => {
    total += getSubTotalPriceNumber(item);
  });
  return total.toLocaleString();
};

export const getSubTotalPrice = (item: OrderItem): string => {
  return getSubTotalPriceNumber(item).toLocaleString();
};

export const getSubTotalPriceNumber = (item: OrderItem): number => {
  return (item.price + getOptionTotalPriceNumber(item.options)) * item.quantity;
};

export const getOptionTotalPrice = (options: OrderOptionItem[]): string => {
  const optTotal = getOptionTotalPriceNumber(options);
  return optTotal.toLocaleString();
};

export const getOptionTotalPriceNumber = (
  options: OrderOptionItem[]
): number => {
  const optTotal = options.reduce(
    (acc: number, current: OrderOptionItem): number => acc + current.price,
    0
  );
  return optTotal;
};

export const getTotalPriceFromCart = (cart: Cart): string => {
  let total = 0;
  Object.keys(cart.items).forEach((key) => {
    const item = cart.items[key];
    total += getSubTotalPriceNumberFromReq(item);
  });
  return total.toLocaleString();
};

export const getTotalBuyItemCount = (cart: Cart): number => {
  let total = 0;
  Object.keys(cart.items).forEach((key) => {
    const item = cart.items[key];
    total += item.quantity;
  });
  return total;
} 

export const getSubTotalPriceFromReq = (req: ItemRequest): string => {
  return getSubTotalPriceNumberFromReq(req).toLocaleString();
};

export const getSubTotalPriceNumberFromReq = (req: ItemRequest): number => {
  const optTotal = getOptionTotalPriceNumber(req.selectOptions);
  return (req.item.price + optTotal) * req.quantity;
};

export const getOptionTotalPriceFromReq = (req: ItemRequest): string => {
  return getOptionTotalPriceNumber(req.selectOptions).toLocaleString();
};
