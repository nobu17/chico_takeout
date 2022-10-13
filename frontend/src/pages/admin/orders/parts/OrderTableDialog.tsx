import * as React from "react";
import Dialog from "@mui/material/Dialog";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import CloseIcon from "@mui/icons-material/Close";
import Slide from "@mui/material/Slide";
import { TransitionProps } from "@mui/material/transitions";
import { AppBar, Toolbar } from "@mui/material";

import { UserOrderInfo } from "../../../../libs/apis/order";
import OrderTable from "../../../../components/parts/OrderTable";
import OrderDetailDialog from "../../../mypage/parts/OrderDetailDialog";

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement;
  },
  ref: React.Ref<unknown>
) {
  return <Slide direction="up" ref={ref} {...props} />;
});

type OrderTableDialogProps = {
  orders: UserOrderInfo[] | undefined;
  title: string;
  open: boolean;
  onClose: () => void;
};

export default function OrderTableDialog(props: OrderTableDialogProps) {
  const [open, setOpen] = React.useState(false);
  const [selectedItem, setSelectedItem] = React.useState<UserOrderInfo>();

  const handleClose = () => {
    props.onClose();
  };

  const handleDetailClose = () => {
    setOpen(false);
  };

  const handleDetailSelect = (item: UserOrderInfo) => {
    setSelectedItem(item);
    setOpen(true);
  };

  if (props.orders === undefined) {
    return <></>;
  }

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
              {props.title}
            </Typography>
          </Toolbar>
        </AppBar>
        <div style={{ height: 900 }}>
          <OrderTable
            orders={props.orders}
            displays={[
              "detailButton",
              "pickupDateTime",
              "userName",
              "userEmail",
              "userTelNo",
              "total",
              "memo",
              "cancel",
              "orderDateTime",
            ]}
            onSelected={handleDetailSelect}
          ></OrderTable>
        </div>
      </Dialog>
      <OrderDetailDialog
        open={open}
        onClose={handleDetailClose}
        order={selectedItem}
      ></OrderDetailDialog>
    </>
  );
}
