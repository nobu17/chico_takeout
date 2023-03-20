import { Button, Stack, Typography } from "@mui/material";
import { ItemInfo, ItemRequest, Cart } from "../../../hooks/UseItemCart";
import ItemSelectTabs from "./ItemSelectTabs";
import { useMessageDialog } from "../../../hooks/UseMessageDialog";

type ItemSelectProps = {
  allItems: CategoryItems[];
  cart: Cart;
  onRequestChanged: callbackRequest;
  onCartUpdated: callbackCartUpdated;
  onSubmit?: callbackSubmit;
  onBack?: callbackBack;
};
interface callbackRequest {
  (item: ItemRequest): void;
}
interface callbackCartUpdated {
  (cart: Cart): void;
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
  const { showMessageDialog, renderDialog } = useMessageDialog();
  const handleSubmit = async () => {
    if (Object.keys(props.cart.items).length <= 0) {
      await showMessageDialog(
        "エラー",
        "注文する商品を１つ以上選択してください。"
      );
      return;
    }
    props.onSubmit?.();
  };
  const handleBack = () => {
    props.onBack?.();
  };
  return (
    <>
      <Typography
        component="h6"
        variant="subtitle2"
        align="center"
        color="error"
        gutterBottom
        sx={{
          mt: 3,
        }}
      >
        ※価格は全て税込みです。<br/>
        上部メニューは→にスライドできます。
      </Typography>
      <ItemSelectTabs
        cart={props.cart}
        allItems={props.allItems}
        onRequestChanged={props.onRequestChanged}
        onCartUpdated={props.onCartUpdated}
      ></ItemSelectTabs>
      <Stack direction="row" spacing={2}>
        <Button
          onClick={handleBack}
          variant="contained"
          color="secondary"
          sx={{ width: 100 }}
        >
          戻る
        </Button>
        <Button onClick={handleSubmit} variant="contained" sx={{ width: 100 }}>
          次へ
        </Button>
      </Stack>
      {renderDialog()}
    </>
  );
}
