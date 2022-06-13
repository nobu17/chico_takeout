import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";

import SpecialHolidayForm from "./SpecialHolidayForm";
import { SpecialHoliday } from "../../../../libs/SpecialHoliday";

type SpecialHolidayFormDialogProps = {
  editItem: SpecialHoliday;
  open: boolean;
  onClose: (item: SpecialHoliday | null) => void;
};

export default function SpecialHolidayFormDialog(
  props: SpecialHolidayFormDialogProps
) {
  const onSubmit = (data: SpecialHoliday) => {
    props.onClose(data);
  };
  const onCancel = () => {
    props.onClose(null);
  };

  return (
    <>
      <Dialog open={props.open} fullWidth maxWidth="sm">
        <DialogContent>
          <SpecialHolidayForm
            editItem={props.editItem}
            onSubmit={onSubmit}
            onCancel={onCancel}
          ></SpecialHolidayForm>
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
}
