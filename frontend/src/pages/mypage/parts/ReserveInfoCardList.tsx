import * as React from "react";
import Typography from "@mui/material/Typography";
import { UserOrderInfo } from "../../../libs/apis/order";
import ReserveInfoCard from "./ReserveInfoCard";
import { OffsetMinutesUserCanCancel } from "../../../libs/Constant";

type ReserveInfoCardListProps = {
  orders?: UserOrderInfo[];
  cancelRequest?: (id: string) => void;
};

export default function ReserveInfoCardList(props: ReserveInfoCardListProps) {
  if (!props.orders || props.orders.length === 0) {
    return (
      <Typography gutterBottom variant="h6" component="div">
        現在、予約はありません。
      </Typography>
    );
  }

  const handleCancel = (id: string) => {
    if (props.cancelRequest && props.orders) {
      props.cancelRequest(id);
    }
  };

  return (
    <>
      <Typography variant="h5">現在の予約</Typography>
      <Typography variant="subtitle1" color="error">※キャンセルは{OffsetMinutesUserCanCancel / 60}時間前まで可能です。</Typography>  
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
