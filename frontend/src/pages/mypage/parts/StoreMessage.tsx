import * as React from "react";
import Typography from "@mui/material/Typography";
import { Stack } from "@mui/material";
import AdbIcon from "@mui/icons-material/Adb";
import { Message } from "../../../libs/Messages";
import LineBreak from "../../../components/parts/LineBreak";

type StoreMessageProps = {
  message?: Message;
};

const getHeader = () => {
  return (
    <Stack
      direction="row"
      justifyContent="center"
      alignItems="center"
      gap={1}
      sx={{ pt: 1, pb: 2 }}
    >
      <AdbIcon></AdbIcon>
      <Typography variant="h5">店舗からのお知らせ</Typography>
    </Stack>
  );
};

export default function StoreMessage(props: StoreMessageProps) {
  if (!props.message) {
    return (
      <>
        {getHeader()}
        <Typography gutterBottom variant="h6" component="div">
          現在、お知らせはありません。
        </Typography>
      </>
    );
  }

  return (
    <>
      {getHeader()}
      <Typography gutterBottom variant="h6" component="div">
        <LineBreak msg={props.message.content}></LineBreak>
      </Typography>
    </>
  );
}
