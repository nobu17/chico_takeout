import * as React from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import { getTotal, UserOrderInfo } from "../../../libs/apis/order";

import { convertDateTimeStrToIncludeDayOfWeeKStr } from "../../../libs/util/DateUtil";
import { Button } from "@mui/material";
import OrderDetailDialog from "./OrderDetailDialog";

type ReserveCardProps = {
  order?: UserOrderInfo;
  cancelRequest?: (id: string) => void;
};

export default function ReserveInfoCard(props: ReserveCardProps) {
  const [open, setOpen] = React.useState(false);

  if (!props.order) {
    return (
      <Typography gutterBottom variant="h6" component="div">
        現在、予約はありません。
      </Typography>
    );
  }

  const handleCancel = () => {
    if (props.cancelRequest && props.order) {
      props.cancelRequest(props.order.id);
    }
  };

  const handleDetailOpen = () => {
    setOpen(true);
  };

  const handleDetailClose = () => {
    setOpen(false);
  };

  return (
    <>
      <Card>
        <CardContent>
          <Typography variant="h6" component="div">
            受取日時:{" "}
            {convertDateTimeStrToIncludeDayOfWeeKStr(
              props.order.pickupDateTime
            )}
          </Typography>
          <Typography color="text.primary">
            合計金額: ¥ {getTotal(props.order).toLocaleString()}
          </Typography>
          <Button
            color="primary"
            variant="contained"
            fullWidth
            onClick={handleDetailOpen}
            sx={{ my: 2 }}
          >
            注文詳細
          </Button>
          <Button
            color="error"
            variant="contained"
            fullWidth
            onClick={handleCancel}
          >
            キャンセルする
          </Button>
        </CardContent>
      </Card>
      <OrderDetailDialog
        open={open}
        onClose={handleDetailClose}
        order={props.order}
      ></OrderDetailDialog>
    </>
  );
}
