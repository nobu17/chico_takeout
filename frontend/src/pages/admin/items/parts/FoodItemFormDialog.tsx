import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";

import FoodItemForm from "./FoodItemForm";
import { ItemKind } from "../../../../libs/ItemKind";
import { FoodItem } from "../../../../libs/FoodItem";
import { BusinessHour } from "../../../../libs/BusinessHour";

type FoodItemFormDialogProps = {
  editItem: FoodItem;
  itemKinds: ItemKind[];
  businessHours: BusinessHour[];
  open: boolean;
  onClose: (item: FoodItem | null) => void;
};

export default function FoodItemFormDialog(props: FoodItemFormDialogProps) {
  const onSubmit = (data: FoodItem) => {
    props.onClose(data);
  };
  const onCancel = () => {
    props.onClose(null);
  };

  return (
    <>
      <Dialog open={props.open} fullWidth maxWidth="sm">
        <DialogContent>
          <FoodItemForm
            itemKinds={props.itemKinds}
            editItem={props.editItem}
            businessHours={props.businessHours}
            onSubmit={onSubmit}
            onCancel={onCancel}
          ></FoodItemForm>
        </DialogContent>
        <DialogActions></DialogActions>
      </Dialog>
    </>
  );
}
