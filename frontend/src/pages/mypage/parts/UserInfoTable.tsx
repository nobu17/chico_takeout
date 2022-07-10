import * as React from "react";
import { UserOrderItem } from "../../../libs/apis/order";
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";

type UserInfoTableProps = {
    userName: string;
    userEmail: string; 
    userTelNo: string;
    memo: string;
};

export default function UserInfoTable(props: UserInfoTableProps) {
  return (
    <TableContainer component={Paper}>
      <Table size="small" aria-label="お客様情報">
        <TableHead>
          <TableRow>
            <TableCell>項目</TableCell>
            <TableCell>入力内容</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          <TableRow>
            <TableCell>氏名</TableCell>
            <TableCell>{props.userName}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>E-mail</TableCell>
            <TableCell>{props.userEmail}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>電話番号</TableCell>
            <TableCell>{props.userTelNo}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>要望やメッセージ等</TableCell>
            <TableCell>
              <Typography
                style={{ display: "inline-block", whiteSpace: "pre-line" }}
              >
                {props.memo}
              </Typography>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </TableContainer>
  );
}
