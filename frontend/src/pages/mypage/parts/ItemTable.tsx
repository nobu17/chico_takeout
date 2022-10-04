import * as React from "react";
import { UserOrderItem } from "../../../libs/apis/order";
import { getTotalByItems } from "../../../libs/apis/order";
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";

type ItemTableProps = {
  foodItems: UserOrderItem[];
  stockItems: UserOrderItem[];
};

export default function ItemTable(props: ItemTableProps) {
  return (
    <TableContainer component={Paper}>
      <Table size="small" aria-label="注文情報">
        <TableHead>
          <TableRow>
            <TableCell>商品名</TableCell>
            <TableCell>個数</TableCell>
            <TableCell>単価</TableCell>
            <TableCell>合計</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {props.foodItems.map((item) => {
            return (
              <TableRow key={item.itemId}>
                <TableCell>{item.name}</TableCell>
                <TableCell>{item.quantity}</TableCell>
                <TableCell>{item.price.toLocaleString()}</TableCell>
                <TableCell>
                  {(item.price * item.quantity).toLocaleString()}
                </TableCell>
              </TableRow>
            );
          })}
          {props.stockItems.map((item) => {
            return (
              <TableRow key={item.itemId}>
                <TableCell>{item.name}</TableCell>
                <TableCell>{item.quantity}</TableCell>
                <TableCell>{item.price.toLocaleString()}</TableCell>
                <TableCell>
                  {(item.price * item.quantity).toLocaleString()}
                </TableCell>
              </TableRow>
            );
          })}
          <TableRow key="summary">
            <TableCell>合計金額</TableCell>
            <TableCell></TableCell>
            <TableCell></TableCell>
            <TableCell>
              {getTotalByItems(
                props.stockItems,
                props.foodItems
              ).toLocaleString()}
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </TableContainer>
  );
}
