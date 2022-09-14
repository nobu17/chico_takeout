import * as React from "react";
import Typography from "@mui/material/Typography";
import { UserOrderInfo } from "../../../libs/apis/order";
import ReserveInfoCard from "./ReserveInfoCard";
import { OffsetMinutesUserCanCancel } from "../../../libs/Constant";
import { Stack } from "@mui/material";
import CalendarMonthIcon from "@mui/icons-material/CalendarMonth";

type ReserveInfoCardListProps = {
  orders?: UserOrderInfo[];
  cancelRequest?: (id: string) => void;
};

const getHeader = () => {
  return (
    <Stack direction="row" justifyContent="center" alignItems="center" gap={1} sx={{ pt: 1, pb: 2 }}>
      <CalendarMonthIcon></CalendarMonthIcon>
      <Typography variant="h5">現在の予約</Typography>
    </Stack>
  );
};

export default function ReserveInfoCardList(props: ReserveInfoCardListProps) {
  if (!props.orders || props.orders.length === 0) {
    return (
      <>
        {getHeader()}
        <Typography gutterBottom variant="h6" component="div">
          現在、予約はありません。
        </Typography>
      </>
    );
  }

  const handleCancel = (id: string) => {
    if (props.cancelRequest && props.orders) {
      props.cancelRequest(id);
    }
  };

  return (
    <>
      {getHeader()}
      <Typography variant="subtitle1" color="error" textAlign="center">
        ※キャンセルは{OffsetMinutesUserCanCancel / 60}時間前まで可能です。
      </Typography>
      {props.orders?.map((order, index) => {
        return (
          <ReserveInfoCard
            key={index}
            order={order}
            cancelRequest={handleCancel}
          ></ReserveInfoCard>
        );
      })}
    </>
  );
}
