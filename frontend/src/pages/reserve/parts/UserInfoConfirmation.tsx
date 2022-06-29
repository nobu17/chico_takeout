import * as React from "react";
import { UserInfo } from "../../../hooks/UseUserInfo";
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

type UserInfoConfirmationProps = {
  userInfo: UserInfo;
};

export default function UserInfoConfirmation(props: UserInfoConfirmationProps) {
  return (
    <TableContainer component={Paper}>
      <Typography gutterBottom variant="h6" align="center" color="text.primary">
        お客様情報
      </Typography>
      <Table aria-label="お客様情報">
        <TableHead>
          <TableRow>
            <TableCell>項目</TableCell>
            <TableCell>入力内容</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          <TableRow>
            <TableCell>氏名</TableCell>
            <TableCell>{props.userInfo.name}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>E-mail</TableCell>
            <TableCell>{props.userInfo.email}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>電話番号</TableCell>
            <TableCell>{props.userInfo.tel}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>要望やメッセージ等</TableCell>
            <TableCell>
              <Typography
                style={{ display: "inline-block", whiteSpace: "pre-line" }}
              >
                {props.userInfo.memo}
              </Typography>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </TableContainer>
  );
}
