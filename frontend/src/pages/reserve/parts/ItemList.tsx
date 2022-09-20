import * as React from "react";
import { Grid } from "@mui/material";
import ItemCard from "./ItemCard";
import Typography from "@mui/material/Typography";
import { Cart, ItemInfo, ItemRequest } from "../../../hooks/UseItemCart";

type ItemListProps = {
  allItems: CategoryItems[];
  cart: Cart;
  onRequestChanged?: callback;
};
interface callback {
  (item: ItemRequest): void;
}

type CategoryItems = {
  title: string;
  items: ItemInfo[];
};

export default function ItemList(props: ItemListProps) {
  const getQuantity = (itemId: string): number => {
    if (props.cart.items[itemId]) {
      return props.cart.items[itemId].quantity;
    }
    return 0;
  };
  return (
    <>
      {props.allItems.map((items, index) => {
        return (
          <div key={index}>
            <Typography
              component="h3"
              variant="h4"
              align="center"
              color="text.primary"
              gutterBottom
              sx={{
                mt: 3,
              }}
            >
              {items.title}
            </Typography>
            <Grid container spacing={2}>
              {items.items.map((item, index) => {
                return (
                  <Grid item xs={12} md={6} key={index + "item"}>
                    <ItemCard
                      item={item}
                      quantity={getQuantity(item.id)}
                      onCountChanged={props.onRequestChanged}
                    ></ItemCard>
                  </Grid>
                );
              })}
            </Grid>
          </div>
        );
      })}
    </>
  );
}
