import * as React from "react";
import { Stack, Button } from "@mui/material";
import { UserInfo } from "../../../hooks/UseUserInfo";
import { Cart } from "../../../hooks/UseItemCart";
import UserInfoConfirmation from "./UserInfoConfirmation";
import ItemConfirmation from "./ItemConfirmation";

type ReserveConfirmationProps = {
  userInfo: UserInfo;
  cart: Cart;
  onSubmit?: callbackSubmit;
  onBack?: callbackBack;
};
interface callbackSubmit {
  (): void;
}

interface callbackBack {
  (): void;
}
export default function ReserveConfirmation(props: ReserveConfirmationProps) {
  return (
    <>
      <UserInfoConfirmation userInfo={props.userInfo} />
      <ItemConfirmation cart={props.cart} />
      <Stack direction="row" spacing={2}>
        <Button variant="contained" onClick={props.onSubmit}>
          注文する
        </Button>
        <Button variant="contained" color="secondary" onClick={props.onBack}>
          戻る
        </Button>
      </Stack>
    </>
  );
}
