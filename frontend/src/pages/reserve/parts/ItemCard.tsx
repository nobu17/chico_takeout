import Counter from "../../../components/parts/Counter";
import * as React from "react";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import CardMedia from "@mui/material/CardMedia";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import{ ItemInfo, ItemRequest } from "../../../hooks/UseItemCart"

type ItemCardProps = {
  item: ItemInfo;
  quantity: number;
  onCountChanged?: callback;
};
interface callback {
  (item :ItemRequest): void;
}

export default function ItemCard(props: ItemCardProps) {
  const onCountChanged = (count: number) => {
    props.onCountChanged?.({ item: props.item, quantity: count });
  };

  return (
    <Card sx={{ maxWidth: 480 }}>
      <CardMedia
        component="img"
        alt="green iguana"
        height="180"
        image="https://mui.com/static/images/cards/contemplative-reptile.jpg"
      />
      <CardContent>
        <Typography gutterBottom variant="h6" component="div">
          {props.item.name}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {props.item.memo}
        </Typography>
      </CardContent>
      <CardActions>
        <Typography sx={{ ml: 1 }} variant="body1">
          ¥ {props.item.price.toLocaleString()}
        </Typography>
        <Box style={{ marginLeft: "auto" }}>
          <Counter
            count={props.quantity}
            max={props.item.max}
            onChanged={onCountChanged}
          ></Counter>
        </Box>
      </CardActions>
    </Card>
  );
}
