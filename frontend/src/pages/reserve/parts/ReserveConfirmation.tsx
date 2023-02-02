import * as React from "react";
import { Stack, Button, Typography } from "@mui/material";
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
      <Typography
        variant="subtitle1"
        align="center"
        color="error"
        gutterBottom
        sx={{
          mt: 3,
        }}
      >
        ※店舗での当日精算になります。
      </Typography>
      <UserInfoConfirmation userInfo={props.userInfo} />
      <ItemConfirmation cart={props.cart} />
      <Stack direction="row" spacing={2}>
        <Button
          variant="contained"
          color="secondary"
          onClick={props.onBack}
          sx={{ width: 100 }}
        >
          戻る
        </Button>
        <Button
          variant="contained"
          onClick={props.onSubmit}
          sx={{ width: 100 }}
        >
          注文する
        </Button>
      </Stack>
    </>
  );
}
