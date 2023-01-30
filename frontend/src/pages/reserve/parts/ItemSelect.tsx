import { Button, Stack, Typography } from "@mui/material";
import { ItemInfo, ItemRequest, Cart } from "../../../hooks/UseItemCart";
import ItemSelectTabs from "./ItemSelectTabs";
import { useMessageDialog } from "../../../hooks/UseMessageDialog";

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
        variant="h6"
        align="center"
        color="error"
        gutterBottom
        sx={{
          mt: 3,
        }}
      >
        ※価格は全て税込みです。
      </Typography>
      <ItemSelectTabs
        cart={props.cart}
        allItems={props.allItems}
        onRequestChanged={props.onRequestChanged}
      ></ItemSelectTabs>
      <Stack direction="row" spacing={2}>
        <Button onClick={handleSubmit} variant="contained">
          次へ
        </Button>
        <Button onClick={handleBack} variant="contained" color="secondary">
          戻る
        </Button>
      </Stack>
      {renderDialog()}
    </>
  );
}
