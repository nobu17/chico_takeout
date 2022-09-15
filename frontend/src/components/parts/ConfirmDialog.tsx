import * as React from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
  Box,
} from "@mui/material";

export type ConfirmDialogProps = {
  open: boolean;
  title: string;
  message: string;
  onClose: (result: boolean) => void;
};

export default function ConfirmDialog(props: ConfirmDialogProps) {
  const { onClose, open } = props;

  const handleOK = () => {
    onClose(true);
  };

  const handleCancel = () => {
    onClose(false);
  };

  return (
    <Dialog open={open}>
      <DialogTitle>{props.title}</DialogTitle>
      <DialogContent>
        <DialogContentText>{props.message}</DialogContentText>
      </DialogContent>
      <DialogActions>
        <Box>
          <Button variant="contained" color="error" onClick={handleCancel}>
            いいえ
          </Button>
          <Button
            sx={{ mx: 2 }}
            variant="contained"
            color="primary"
            onClick={handleOK}
          >
            はい
          </Button>
        </Box>
      </DialogActions>
    </Dialog>
  );
}
