import {
  MenuItem,
  TextField,
  TextFieldProps,
  ListSubheader,
} from "@mui/material";
import type { ChangeEventHandler, FocusEventHandler } from "react";

export type GroupSelectsProps = {
  label: string;
  multiple: boolean;
  itemList: CategoryItems[];
  error?: string;
};

type CategoryItems = {
  groupName: string;
  value: string[];
};

type FlatItem = {
  isHeader: boolean;
  value: string;
};

const convertFlat = (categories: CategoryItems[]): FlatItem[] => {
  const flats: FlatItem[] = [];
  for (const category of categories) {
    flats.push({ isHeader: true, value: category.groupName });
    for (const val of category.value) {
      flats.push({ isHeader: false, value: val });
    }
  }
  return flats;
};

const convertMenuItem = (categories: CategoryItems[]) => {
  const flats = convertFlat(categories);
  return flats.map((f) =>
    f.isHeader ? (
      <ListSubheader key={f.value}>{f.value}</ListSubheader>
    ) : (
      <MenuItem key={f.value} value={f.value}>
        {f.value}
      </MenuItem>
    )
  );
};

export function GroupSelects(
  props: GroupSelectsProps & {
    inputRef: TextFieldProps["ref"];
    value: string;
    onChange: ChangeEventHandler<HTMLTextAreaElement>;
    onBlur: FocusEventHandler<HTMLTextAreaElement>;
  }
) {
  return (
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
      {convertMenuItem(props.itemList)}
    </TextField>
  );
}
