import * as React from "react";
import { Cart, OptionItemInfo } from "../../../hooks/UseItemCart";
import OrderDetailTable from "../../../components/parts/OrderDetailTable";

type ItemConfirmationProps = {
  cart: Cart;
};

type OrderItem = {
  itemId: string;
  name: string;
  price: number;
  quantity: number;
  options: OrderOptionItem[];
};

type OrderOptionItem = {
  itemId: string;
  name: string;
  price: number;
};

const convertOrderItem = (cart: Cart): OrderItem[] => {
  let items: OrderItem[] = [];
  Object.keys(cart.items).forEach((key) => {
    const cartItem = cart.items[key];
    items.push({
      itemId: cartItem.item.id,
      name: cartItem.item.name,
      price: cartItem.item.price,
      quantity: cartItem.quantity,
      options: convertOptionItem(cartItem.selectOptions),
    });
  });
  return items;
};

const convertOptionItem = (items: OptionItemInfo[]): OrderOptionItem[] => {
  return items.map((item) => {
    return {
      itemId: item.id,
      name: item.name,
      price: item.price,
    };
  });
};

export default function ItemConfirmation(props: ItemConfirmationProps) {
  return (
    <OrderDetailTable items={convertOrderItem(props.cart)}></OrderDetailTable>
  );
}
