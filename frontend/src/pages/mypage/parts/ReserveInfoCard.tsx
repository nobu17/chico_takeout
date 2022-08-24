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

import { convertDateTimeStrToIncludeDayOfWeeKStr } from "../../../libs/util/DateUtil";
import { Button } from "@mui/material";

type ReserveCardProps = {
  order?: UserOrderInfo;
  cancelRequest?: (id: string) => void;
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
    </>
  );
}
