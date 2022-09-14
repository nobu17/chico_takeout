import * as React from "react";
import Dialog from "@mui/material/Dialog";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import CloseIcon from "@mui/icons-material/Close";
import Slide from "@mui/material/Slide";
import { TransitionProps } from "@mui/material/transitions";
import { getTotal, UserOrderInfo } from "../../../libs/apis/order";
import { AppBar, Toolbar } from "@mui/material";

import ItemTable from "./ItemTable";
import UserInfoTable from "./UserInfoTable";
import { convertDateTimeStrToIncludeDayOfWeeKStr } from "../../../libs/util/DateUtil";

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement;
  },
  ref: React.Ref<unknown>
) {
  return <Slide direction="up" ref={ref} {...props} />;
});

type OrderDetailDialogProps = {
  order: UserOrderInfo;
  open: boolean;
  onClose: () => void;
};

export default function OrderDetailDialog(props: OrderDetailDialogProps) {
  const handleClose = () => {
    props.onClose();
  };

  return (
    <>
      <Dialog
        open={props.open}
        fullScreen
        onClose={handleClose}
        TransitionComponent={Transition}
      >
        <AppBar sx={{ position: "relative" }}>
          <Toolbar>
            <IconButton
              edge="start"
              color="inherit"
              onClick={handleClose}
              aria-label="close"
            >
              <CloseIcon />
            </IconButton>
            <Typography sx={{ ml: 2, flex: 1 }} variant="h6" component="div">
              注文詳細
            </Typography>
          </Toolbar>
        </AppBar>
        <Typography variant="h6" sx={{ px: 1, py: 2 }}>
          受取日時:{" "}
          {convertDateTimeStrToIncludeDayOfWeeKStr(props.order.pickupDateTime)}
        </Typography>
        <Typography textAlign="center" variant="h5" sx={{ py: 2 }}>
          注文商品
        </Typography>
        <ItemTable {...props.order}></ItemTable>
        <Typography variant="h6" sx={{ px: 2, py: 1 }} textAlign="right">
          合計金額: ¥ {getTotal(props.order).toLocaleString()}
        </Typography>
        <Typography textAlign="center" variant="h5" sx={{ py: 2 }}>
          お客様情報
        </Typography>
        <UserInfoTable {...props.order}></UserInfoTable>
      </Dialog>
    </>
  );
}
