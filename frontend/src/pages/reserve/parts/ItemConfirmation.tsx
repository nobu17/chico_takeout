import * as React from "react";
import { Cart } from "../../../hooks/UseItemCart";
import {
  Table,
  TableContainer,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Paper,
  Typography,
} from "@mui/material";

type ItemConfirmationProps = {
  cart: Cart;
};

export default function ItemConfirmation(props: ItemConfirmationProps) {
  const getTotalPrice = (cart: Cart): string => {
    let total = 0;
    Object.keys(props.cart.items).forEach((key, index) => {
      total +=
        props.cart.items[key].item.price * props.cart.items[key].quantity;
    });
    return total.toLocaleString();
  };
  return (
    <TableContainer component={Paper}>
      <Typography gutterBottom variant="h6" align="center" color="text.primary">
        注文情報
      </Typography>
      <Table aria-label="注文情報">
        <TableHead>
          <TableRow>
            <TableCell>商品名</TableCell>
            <TableCell>個数</TableCell>
            <TableCell>価格</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {Object.keys(props.cart.items).map((key) => {
            return (
              <TableRow key={key}>
                <TableCell>{props.cart.items[key].item.name}</TableCell>
                <TableCell>{props.cart.items[key].quantity}</TableCell>
                <TableCell>
                  {(
                    props.cart.items[key].item.price *
                    props.cart.items[key].quantity
                  ).toLocaleString()}
                </TableCell>
              </TableRow>
            );
          })}
          <TableRow>
            <TableCell>合計</TableCell>
            <TableCell></TableCell>
            <TableCell>¥ {getTotalPrice(props.cart)}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </TableContainer>
  );
}
