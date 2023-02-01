import * as React from "react";
import FloatingCartIconButton from "../../../components/parts/FloatingCartIconButton";
import CartEditDialog from "./CartEditDialog";
import { getTotalBuyItemCount } from "../../../libs/util/ItemCalc";
import { Cart, ItemInfo } from "../../../hooks/UseItemCart";

type CartButtonProps = {
  allItems: CategoryItems[];
  cart: Cart;
  onUpdated: (cart: Cart) => void;
};

type CategoryItems = {
  title: string;
  items: ItemInfo[];
};

export default function CartButton(props: CartButtonProps) {
  const [open, setOpen] = React.useState(false);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleCancel = () => {
    setOpen(false);
  };
  const handleSubmit = (cart: Cart) => {
    setOpen(false);
    props.onUpdated(cart);
  };

  // only when cart has item, display buttons
  if (
    !props.cart ||
    !props.cart.items ||
    Object.keys(props.cart.items).length === 0
  ) {
    return <></>;
  }
  return (
    <>
      <FloatingCartIconButton
        onClick={handleOpen}
        count={getTotalBuyItemCount(props.cart)}
      ></FloatingCartIconButton>
      <CartEditDialog
        open={open}
        cart={props.cart}
        allItems={props.allItems}
        onCancel={handleCancel}
        onSubmit={handleSubmit}
      ></CartEditDialog>
    </>
  );
}
