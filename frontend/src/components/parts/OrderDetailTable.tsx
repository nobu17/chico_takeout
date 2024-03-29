import * as React from "react";
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

import { getTotalPrice, getSubTotalPrice } from "../../libs/util/ItemCalc";

type OrderDetailTableProps = {
  items: OrderItem[];
};

type OrderItem = {
  itemId: string;
  name: string;
  price: number;
  quantity: number;
  options: OrderOptionItem[];
};

type OrderOptionItem = {
  itemId: string;
  name: string;
  price: number;
};

export default function OrderDetailTable(props: OrderDetailTableProps) {
  return (
    <TableContainer component={Paper}>
      <Typography textAlign="center" variant="h5" sx={{ py: 2 }}>
        注文商品
      </Typography>
      <Table aria-label="注文情報" size="small">
        <TableHead>
          <TableRow>
            <TableCell style={{ width: "47%" }}>商品名</TableCell>
            <TableCell style={{ width: "25%" }} align="right">
              個数
            </TableCell>
            <TableCell style={{ width: "28%" }} align="right">
              価格
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {props.items.map((item) => {
            return (
              <React.Fragment key={item.itemId}>
                <TableRow>
                  <TableCell>{item.name}</TableCell>
                  <TableCell align="right">{item.quantity}</TableCell>
                  <TableCell align="right">
                    {item.price.toLocaleString()}
                  </TableCell>
                </TableRow>
                {item.options.map((opt) => {
                  return (
                    <TableRow key={opt.itemId}>
                      <TableCell>{`(${opt.name})`}</TableCell>
                      <TableCell></TableCell>
                      <TableCell align="right">
                        {opt.price.toLocaleString()}
                      </TableCell>
                    </TableRow>
                  );
                })}
                <TableRow>
                  <TableCell />
                  <TableCell>小計</TableCell>
                  <TableCell align="right">
                    ¥ {getSubTotalPrice(item)}
                  </TableCell>
                </TableRow>
              </React.Fragment>
            );
          })}
          <TableRow style={{ borderTop: "dotted" }}>
            <TableCell></TableCell>
            <TableCell>合計</TableCell>
            <TableCell align="right">¥ {getTotalPrice(props.items)}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </TableContainer>
  );
}
