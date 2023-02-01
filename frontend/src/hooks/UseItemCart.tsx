import { useState } from "react";

export type Cart = {
  items: { [id: string]: ItemRequest };
};

export type ItemRequest = {
  item: ItemInfo;
  quantity: number;
  selectOptions: OptionItemInfo[];
};

export type ItemInfo = {
  id: string;
  type: "food" | "stock";
  name: string;
  imageUrl: string;
  memo: string;
  price: number;
  max: number;
  options: OptionItemInfo[];
};

export type OptionItemInfo = {
  id: string;
  name: string;
  description: string;
  price: number;
};

const defaultCart: Cart = {
  items: {},
};

export function useItemCart(init: Cart = defaultCart) {
  const [cart, setCart] = useState<Cart>(JSON.parse(JSON.stringify(init)));
  const updateCart = (request: ItemRequest) => {
    const newCarts: Cart = { ...cart };
    if (request.quantity > 0) {
      newCarts.items[request.item.id] = request;
    } else {
      delete newCarts.items[request.item.id];
    }
    setCart(newCarts);
  };

  // this method allow quantity 0. for draft usage
  const updateCartAsDraft = (request: ItemRequest) => {
    const newCarts: Cart = { ...cart };
    newCarts.items[request.item.id] = request;
    setCart(newCarts);
  };

  // remove quantity zero item
  const getActivatedCartFromDraft = (): Cart => {
    const newCart: Cart = { ...cart };
    Object.keys(cart.items).forEach((key) => {
      if (cart.items[key].quantity <= 0) {
        delete newCart.items[key];
      }
    });
    return newCart;
  };

  const resetCart = (init: Cart = defaultCart) => {
    setCart(JSON.parse(JSON.stringify(init)));
  };

  return {
    cart,
    updateCart,
    updateCartAsDraft,
    getActivatedCartFromDraft,
    resetCart,
  };
}
