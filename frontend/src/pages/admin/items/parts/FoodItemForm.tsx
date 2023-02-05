import * as React from "react";
import {
  Container,
  MenuItem,
  Stack,
  TextField,
  FormControlLabel,
  Checkbox,
} from "@mui/material";

import { ItemKind } from "../../../../libs/ItemKind";
import { FoodItem } from "../../../../libs/FoodItem";
import { BusinessHour } from "../../../../libs/BusinessHour";
import { RhfSelects } from "../../../../components/parts/Rhf/RhfSelects";
import { SubmitHandler, useForm, Controller } from "react-hook-form";
import SubmitButtons from "../../../../components/parts/SubmitButtons";
import { FirebaseImageUpload } from "../../../../components/parts/FirebaseImageUpload";

import {
  RequiredErrorMessage,
  MaxLengthErrorMessage,
  MaxErrorMessage,
  MinErrorMessage,
} from "../../../../libs/ErrorMessages";

import MultiSelectDatePicker from "../../../../components/parts/MultiSelectDatePicker";

const baseUrl = "/foods";

type FoodItemFormProps = {
  editItem: FoodItem;
  itemKinds: ItemKind[];
  businessHours: BusinessHour[];
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (foodItem: FoodItem): void;
}
interface callbackCancel {
  (): void;
}

export type FoodItemInput = {
  id: string;
  name: string;
  description: string;
  priority: number;
  price: number;
  maxOrder: number;
  maxOrderPerDay: number;
  enabled: boolean;
  imageUrl: string;
  kindId: string; // for select only using kindId
  scheduleIds: string[];
  allowDates: string[];
};

type HourSelects = {
  name: string;
  value: string;
};

const convertInput = (item: FoodItem): FoodItemInput => {
  const kindId = item.kind ? item.kind.id : "";
  return {
    id: item.id,
    name: item.name,
    description: item.description,
    priority: item.priority,
    price: item.price,
    maxOrder: item.maxOrder,
    maxOrderPerDay: item.maxOrderPerDay,
    enabled: item.enabled,
    imageUrl: item.imageUrl,
    kindId: kindId,
    scheduleIds: item.scheduleIds,
    allowDates: item.allowDates,
  };
};

const reverseInput = (item: FoodItemInput, kinds: ItemKind[]): FoodItem => {
  const kind = kinds.find((k) => k.id === item.kindId);
  return {
    id: item.id,
    name: item.name,
    description: item.description,
    priority: item.priority,
    price: item.price,
    maxOrder: item.maxOrder,
    maxOrderPerDay: item.maxOrderPerDay,
    enabled: item.enabled,
    imageUrl: item.imageUrl,
    kind: kind!,
    scheduleIds: item.scheduleIds,
    allowDates: item.allowDates,
  };
};

const convertHoursList = (hours: BusinessHour[]): HourSelects[] => {
  return hours.map((val) => ({
    name: val.name,
    value: val.id,
  }));
};

export default function FoodItemForm(props: FoodItemFormProps) {
  const {
    register,
    setValue,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<FoodItemInput>({ defaultValues: convertInput(props.editItem) });

  const onSubmit: SubmitHandler<FoodItemInput> = (data) => {
    const converted = reverseInput(data, props.itemKinds);
    props.onSubmit(converted);
  };

  const onCancel = () => {
    props.onCancel();
  };

  const onImageChanged = (fileUrl: string) => {
    setValue("imageUrl", fileUrl);
  };

  const hoursList = convertHoursList(props.businessHours);

  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 2 }}>
        <Stack spacing={3}>
          <TextField
            label="名称"
            {...register("name", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 25, message: MaxLengthErrorMessage(25) },
            })}
            error={Boolean(errors.name)}
            helperText={errors.name && errors.name.message}
          />
          <TextField
            label="説明"
            multiline
            rows={5}
            {...register("description", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 500, message: MaxLengthErrorMessage(500) },
            })}
            error={Boolean(errors.description)}
            helperText={errors.description && errors.description.message}
          />
          <Controller
            name="kindId"
            control={control}
            defaultValue={""}
            rules={{ required: true }}
            render={({ field: { onChange, value, ref } }) => (
              <TextField
                inputRef={ref}
                onChange={onChange}
                value={value}
                select
                sx={{ mt: 2 }}
                fullWidth
                label="アイテム種別"
                error={Boolean(errors.kindId)}
                helperText={errors.kindId && errors.kindId.message}
              >
                {props.itemKinds.map((kind, index) => (
                  <MenuItem key={index} value={kind.id}>
                    {kind.name}
                  </MenuItem>
                ))}
              </TextField>
            )}
          />
          <TextField
            type="number"
            label="表示順序"
            {...register("priority", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 30, message: MaxErrorMessage(30) },
            })}
            error={Boolean(errors.priority)}
            helperText={errors.priority && errors.priority.message}
          />
          <TextField
            type="number"
            label="価格"
            {...register("price", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 20000, message: MaxErrorMessage(20000) },
            })}
            error={Boolean(errors.price)}
            helperText={errors.price && errors.price.message}
          />
          <TextField
            type="number"
            label="最大注文可能数"
            {...register("maxOrder", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 30, message: MaxErrorMessage(30) },
            })}
            error={Boolean(errors.maxOrder)}
            helperText={errors.maxOrder && errors.maxOrder.message}
          />
          <TextField
            type="number"
            label="在庫数(日別)"
            {...register("maxOrderPerDay", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 100, message: MaxErrorMessage(100) },
            })}
            error={Boolean(errors.maxOrderPerDay)}
            helperText={errors.maxOrderPerDay && errors.maxOrderPerDay.message}
          />
          <RhfSelects
            label="販売時間"
            name="scheduleIds"
            multiple={true}
            itemList={hoursList}
            control={control}
          />
          <Controller
            name="allowDates"
            control={control}
            render={({ field: { onChange, value, ref } }) => (
              <MultiSelectDatePicker
                selectedDates={value}
                onCancel={() => {}}
                onSubmit={(val) => { setValue("allowDates", val) }}
              />
            )}
          />
          <Controller
            name="imageUrl"
            control={control}
            render={({ field: { onChange, value, ref } }) => (
              <FirebaseImageUpload
                baseUrl={baseUrl}
                imageUrl={value}
                onImageUploaded={onImageChanged}
              ></FirebaseImageUpload>
            )}
          />
          <FormControlLabel
            control={
              <Controller
                defaultValue={true}
                name="enabled"
                control={control}
                render={({ field: { onChange, value, ref } }) => (
                  <Checkbox
                    inputRef={ref}
                    checked={value}
                    onChange={onChange}
                  />
                )}
              />
            }
            label="有効"
          />
          <SubmitButtons
            onSubmit={handleSubmit(onSubmit)}
            onCancel={onCancel}
          />
        </Stack>
      </Container>
    </>
  );
}
