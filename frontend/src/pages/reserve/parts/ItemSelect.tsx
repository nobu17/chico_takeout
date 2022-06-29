import { Button, Stack } from "@mui/material";
import { ItemInfo, ItemRequest, Cart } from "../../../hooks/UseItemCart";
import ItemList from "./ItemList";

type ItemSelectProps = {
  allItems: CategoryItems[];
  cart: Cart;
  onRequestChanged?: callbackRequest;
  onSubmit?: callbackSubmit;
  onBack?: callbackBack;
};
interface callbackRequest {
  (item: ItemRequest): void;
}
interface callbackSubmit {
  (): void;
}
interface callbackBack {
  (): void;
}

type CategoryItems = {
  title: string;
  items: ItemInfo[];
};

export default function ItemSelect(props: ItemSelectProps) {
  const handleSubmit = () => {
    if (Object.keys(props.cart.items).length <= 0) {
      alert("注文するアイテムが選択されていません。");
      return;
    }
    props.onSubmit?.();
  };
  const handleBack = () => {
    props.onBack?.();
  };

  return (
    <>
      <ItemList
        cart={props.cart}
        allItems={props.allItems}
        onRequestChanged={props.onRequestChanged}
      ></ItemList>
      <Stack direction="row" spacing={2}>
        <Button onClick={handleSubmit} variant="contained">
          次へ
        </Button>
        <Button onClick={handleBack} variant="contained" color="secondary">
          戻る
        </Button>
      </Stack>
    </>
  );
}
