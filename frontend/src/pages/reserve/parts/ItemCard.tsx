import * as React from "react";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import CardMedia from "@mui/material/CardMedia";
import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";

import LineBreak from "../../../components/parts/LineBreak";
import Counter from "../../../components/parts/Counter";
import OptionItemSelectDialog from "./OptionItemSelectDialog";
import {
  ItemInfo,
  ItemRequest,
  OptionItemInfo,
} from "../../../hooks/UseItemCart";
import { getImageUrl } from "../../../libs/util/ImageUtil";

type ItemCardProps = {
  item: ItemInfo;
  quantity: number;
  selectedOptions: OptionItemInfo[];
  onChanged?: callback;
};
interface callback {
  (item: ItemRequest): void;
}

const getOptionCountString = (selected: OptionItemInfo[]): string => {
  if (selected.length < 1) {
    return "";
  }
  return `(${selected.length} 個)`;
};

const getOptionTotal = (options: OptionItemInfo[]): string => {
  const totalStr = options
    .reduce(
      (acc: number, current: OptionItemInfo): number => acc + current.price,
      0
    )
    .toLocaleString();

  return `  (+¥ ${totalStr})`;
};

export default function ItemCard(props: ItemCardProps) {
  const [open, setOpen] = React.useState(false);
  const [options, setOptions] = React.useState(props.selectedOptions);

  // need update from cart button update
  React.useEffect(() => {
    setOptions(props.selectedOptions);
  },[props.selectedOptions])

  const onCountChanged = (count: number) => {
    props.onChanged?.({
      item: props.item,
      quantity: count,
      selectOptions: options,
    });
  };

  const handleOptionClick = () => {
    setOpen(true);
  };

  const handleOptionClose = (selected: OptionItemInfo[]) => {
    setOptions(selected);
    setOpen(false);
    props.onChanged?.({
      item: props.item,
      quantity: props.quantity,
      selectOptions: selected,
    });
  };

  return (
    <Card sx={{ maxWidth: 480 }}>
      <CardMedia
        component="img"
        height="230"
        sx={{ objectFit: "contain" }}
        image={getImageUrl(props.item.imageUrl)}
      />
      <CardContent>
        <Typography gutterBottom variant="h6" component="div">
          {props.item.name}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          <LineBreak msg={props.item.memo} />
        </Typography>
      </CardContent>
      <CardActions>
        <Stack>
          <Typography sx={{ ml: 1 }} variant="body1">
            {" "}
            ¥ {props.item.price.toLocaleString()}
          </Typography>
          {options.length < 1 ? (
            <></>
          ) : (
            <Typography sx={{ ml: 1 }} variant="body2">
              {" "}
              {getOptionTotal(options)}
            </Typography>
          )}
        </Stack>

        <Box style={{ marginLeft: "auto" }}>
          <Counter
            count={props.quantity}
            max={props.item.max}
            onChanged={onCountChanged}
          ></Counter>
        </Box>
      </CardActions>
      {props.item.options.length <= 0 ? (
        <></>
      ) : (
        <Box sx={{ m: 2 }}>
          <Button
            variant="contained"
            color="error"
            fullWidth
            onClick={() => handleOptionClick()}
          >
            オプション {getOptionCountString(options)}
          </Button>
        </Box>
      )}
      <OptionItemSelectDialog
        open={open}
        items={props.item.options}
        selected={options}
        onClose={handleOptionClose}
      />
    </Card>
  );
}
