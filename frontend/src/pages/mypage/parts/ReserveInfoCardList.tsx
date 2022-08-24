import * as React from "react";
import Typography from "@mui/material/Typography";
import { UserOrderInfo } from "../../../libs/apis/order";
import ReserveInfoCard from "./ReserveInfoCard";

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
      {props.orders?.map((order) => {
        return (
          <ReserveInfoCard
            order={order}
            cancelRequest={handleCancel}
          ></ReserveInfoCard>
        );
      })}
    </>
  );
}
