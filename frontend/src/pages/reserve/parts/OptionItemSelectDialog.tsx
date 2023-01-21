import * as React from "react";
import Dialog from "@mui/material/Dialog";
import { TransitionProps } from "@mui/material/transitions";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Checkbox from "@mui/material/Checkbox";
import { Typography, AppBar, Toolbar, IconButton, Slide } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { OptionItemInfo } from "../../../hooks/UseItemCart";

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement;
  },
  ref: React.Ref<unknown>
) {
  return <Slide direction="up" ref={ref} {...props} />;
});

type OptionItemSelectDialogProps = {
  items: OptionItemInfo[];
  selected: OptionItemInfo[];
  open: boolean;
  onClose: (items: OptionItemInfo[]) => void;
};

type OptionItemInput = {
  item: OptionItemInfo;
  checked: boolean;
};

const convertInput = (items: OptionItemInfo[], selected: OptionItemInfo[]): OptionItemInput[] => {
  const converted: OptionItemInput[] = [];
  for (const item of items) {
    converted.push({
      item: item,
      checked: selected.some(s => s.id === item.id),
    });
  }
  return converted;
};

export default function OptionItemSelectDialog(
  props: OptionItemSelectDialogProps
) {
  const [inputs, setInputs] = React.useState(convertInput(props.items, props.selected));

  const handleToggle = (item: OptionItemInfo) => () => {
    let targetItem = inputs.find((x) => x.item.id === item.id);
    if (targetItem) {
      targetItem.checked = !targetItem.checked;
    }
    const newInputs = [...inputs];
    setInputs(newInputs);
  };

  const handleClose = () => {
    const selected = inputs.filter((f) => f.checked).map((f) => f.item);
    props.onClose(selected);
  };
  return (
    <>
      <Dialog
        open={props.open}
        fullScreen
        onClose={handleClose}
        TransitionComponent={Transition}
      >
        <AppBar sx={{ position: "relative" }}>
          <Toolbar>
            <IconButton
              edge="start"
              color="inherit"
              onClick={handleClose}
              aria-label="close"
            >
              <CloseIcon />
            </IconButton>
            <Typography sx={{ ml: 2, flex: 1 }} variant="h6" component="div">
              オプション選択
            </Typography>
          </Toolbar>
        </AppBar>
        <List
          sx={{ width: "100%", maxWidth: 360, bgcolor: "background.paper" }}
        >
          {inputs.map((opt) => {
            const item = opt.item;
            const labelId = `checkbox-list-label-${item.id}`;

            return (
              <ListItem key={item.id} disablePadding>
                <ListItemButton
                  role={undefined}
                  onClick={handleToggle(item)}
                  dense
                >
                  <ListItemIcon>
                    <Checkbox
                      edge="start"
                      checked={opt.checked}
                      tabIndex={-1}
                      disableRipple
                      inputProps={{ "aria-labelledby": labelId }}
                    />
                  </ListItemIcon>
                  <ListItemText
                    id={labelId}
                    primary={`${item.name} : ¥ ${item.price.toLocaleString()}`}
                    secondary={
                      <React.Fragment>{item.description}</React.Fragment>
                    }
                  />
                </ListItemButton>
              </ListItem>
            );
          })}
        </List>
      </Dialog>
    </>
  );
}
