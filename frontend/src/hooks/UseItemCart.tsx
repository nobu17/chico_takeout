import { useState } from "react";

export type Cart = {
  items: { [id: string]: ItemRequest };
};

export type ItemRequest = {
  item: ItemInfo;
  quantity: number;
};

export type ItemInfo = {
  id: string;
  name: string;
  imageUrl: string;
  memo: string;
  price: number;
  max: number;
};

const defaultCart: Cart = {
  items: {},
};

export function useItemCart() {
  const [cart, setCart] = useState<Cart>(JSON.parse(JSON.stringify(defaultCart)));
  const updateCart = (request: ItemRequest) => {
    const newCarts: Cart = {...cart};
    if (request.quantity > 0) {
      newCarts.items[request.item.id] = request;
    } else {
      delete newCarts.items[request.item.id];
    }
    setCart(newCarts);
  };

  const resetCart = () => {
    setCart(JSON.parse(JSON.stringify(defaultCart)));
  };

  return {
    cart,
    updateCart,
    resetCart,
  };
}
