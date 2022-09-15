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

export type MessageDialogProps = {
  open: boolean;
  title: string;
  message: string;
  onClose: () => void;
};

export default function MessageDialog(props: MessageDialogProps) {
  const { onClose, open } = props;

  const handleClose = () => {
    onClose();
  };

  return (
    <Dialog onClose={handleClose} open={open}>
      <DialogTitle>{props.title}</DialogTitle>
      <DialogContent>
        <DialogContentText>{props.message}</DialogContentText>
      </DialogContent>
      <DialogActions>
        <Box>
          <Button onClick={handleClose}>OK</Button>
        </Box>
      </DialogActions>
    </Dialog>
  );
}
