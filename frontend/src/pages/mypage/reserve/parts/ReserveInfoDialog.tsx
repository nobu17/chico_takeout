import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import { Button, Typography } from "@mui/material";
import { UserOrderInfo } from "../../../../libs/apis/order";
import UserInfoTable from "../../parts/UserInfoTable";
import OrderDetailTable from "../../../../components/parts/OrderDetailTable";

type ReserveInfoDialogProps = {
  item?: UserOrderInfo;
  open: boolean;
  onClose: () => void;
};

export default function ReserveInfoDialog(props: ReserveInfoDialogProps) {
  if (!props.item) {
    return <></>;
  }
  return (
    <>
      <Dialog open={props.open} onClose={props.onClose} fullWidth maxWidth="sm">
        <DialogContent>
          <Typography sx={{ m: 2 }}>商品情報</Typography>
          <OrderDetailTable
            items={props.item.foodItems.concat(props.item.stockItems)}
          ></OrderDetailTable>
          <Typography sx={{ m: 2 }}>お客様情報</Typography>
          <UserInfoTable {...props.item}></UserInfoTable>
        </DialogContent>
        <DialogActions>
          <Button
            color="primary"
            variant="contained"
            size="large"
            onClick={props.onClose}
          >
            閉じる
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
