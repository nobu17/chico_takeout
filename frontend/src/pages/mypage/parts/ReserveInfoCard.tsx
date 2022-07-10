import * as React from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";

import Typography from "@mui/material/Typography";
import { UserOrderInfo } from "../../../libs/apis/order";
import ItemTable from "./ItemTable";
import UserInfoTable from "./UserInfoTable";

type ReserveCardProps = {
  order?: UserOrderInfo;
};

const getTotal = (order: UserOrderInfo): number => {
  const stockTotal = order.stockItems.reduce((sum, element) => {
    return sum + element.price * element.quantity;
  }, 0);
  const foodTotal = order.foodItems.reduce((sum, element) => {
    return sum + element.price * element.quantity;
  }, 0);
  return stockTotal + foodTotal;
};

export default function ReserveInfoCard(props: ReserveCardProps) {
  console.log("ReserveInfoCard", props);
  if (!props.order) {
    return (
      <Typography gutterBottom variant="h6" component="div">
        現在、予約はありません。
      </Typography>
    );
  }
  return (
    <>
      <Typography variant="h5">
        現在の予約
      </Typography>
      <Card>
        <CardContent>
          <Typography variant="h6" component="div">
            受取日時: 2020/02/01(金) 15:09
          </Typography>
          <Typography color="text.primary">
            合計金額: ¥ {getTotal(props.order).toLocaleString()}
          </Typography>
          <Accordion sx={{ my: 2 }}>
            <AccordionSummary
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel1a-content"
            >
              <Typography>商品情報</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <ItemTable {...props.order}></ItemTable>
            </AccordionDetails>
          </Accordion>
          <Accordion sx={{ my: 2 }}>
            <AccordionSummary
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel1a-content"
            >
              <Typography>お客様情報</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <UserInfoTable {...props.order}></UserInfoTable>
            </AccordionDetails>
          </Accordion>
        </CardContent>
      </Card>
    </>
  );
}
