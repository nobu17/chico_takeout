import { MenuItem, TextField, TextFieldProps } from "@mui/material";
import type { ChangeEventHandler, FocusEventHandler } from "react";

export type SelectsProps = {
  label: string;
  multiple: boolean;
  itemList: SelectValue[];
  error?: string;
};

type SelectValue = {
    name: string
    value: string
}

export function Selects(
  props: SelectsProps & {
    inputRef: TextFieldProps["ref"];
    value: string[];
    onChange: ChangeEventHandler<HTMLTextAreaElement>;
    onBlur: FocusEventHandler<HTMLTextAreaElement>;
  }
) {
  return (
    <>
      <TextField
        ref={props.inputRef}
        onChange={props.onChange}
        onBlur={props.onBlur}
        select
        SelectProps={{
          multiple: props.multiple,
          value: props.value,
        }}
        sx={{ mt: 2 }}
        fullWidth
        label={props.label}
        error={Boolean(props.error)}
        helperText={props.error}
      >
        {props.itemList.map((item, index) => (
          <MenuItem key={index} value={item.value}>
            {item.name}
          </MenuItem>
        ))}
      </TextField>
    </>
  );
}
