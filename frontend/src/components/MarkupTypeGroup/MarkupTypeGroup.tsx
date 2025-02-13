import { Checkbox, Radio, Select, TextInput } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import { MarkupTypeField } from "../../utils/types";
import "./MarkupTypeGroup.scss";
import { useMemo } from "react";

const _ = require("lodash");

const b = block("markup-type-group");

export const MarkupTypeGroup = ({
  fields,
  onUpdate,
  value,
}: {
  fields: MarkupTypeField[];
  onUpdate: (value: string[]) => void;
  value?: string[];
}) => {
  const uniqueId = useMemo(() => _.uniqueId(), []);

  if (fields.length === 0 || value === undefined) {
    return <></>;
  }
  if (fields[0].assessment_type_id === 1) {
    return (
      <div className={b("radio")}>
        <label className={b("label")}>{fields[0].label}</label>
        {fields.map((field) => (
          <Radio
            id={`${uniqueId}_${String(field.group_id)}_${String(field.id)}`}
            value={String(field.id)}
            content={field.name}
            onUpdate={(checked) => {
              let value: string[] = [];
              if (checked) {
                value = [String(field.id)];
              }
              onUpdate(value);
            }}
            checked={value[0] === String(field.id)}
          />
        ))}
      </div>
    );
  } else if (fields[0].assessment_type_id === 2) {
    return (
      <div className={b("checkbox")}>
        <label className={b("label")}>{fields[0].label}</label>
        {fields.map((field) => (
          <Checkbox
            id={`${uniqueId}_${String(field.group_id)}_${String(field.id)}`}
            value={String(field.id)}
            content={field.name}
            onUpdate={(checked) => {
              let valueCopy: string[] = _.cloneDeep(value);
              if (checked) {
                if (!valueCopy.includes(String(field.id))) {
                  valueCopy.push(String(field.id));
                }
              } else {
                valueCopy.splice(valueCopy.indexOf(String(field.id)), 1);
              }
              onUpdate(valueCopy);
            }}
            checked={value.includes(String(field.id))}
          />
        ))}
      </div>
    );
  } else if (fields[0].assessment_type_id === 3) {
    return (
      <div className={b("select")}>
        <label
          htmlFor={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          className={b("label")}
        >
          {fields[0].label}
        </label>
        <Select
          placeholder={fields[0].label}
          width={"max"}
          id={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          value={value}
          onUpdate={onUpdate}
        >
          {fields.map((field) => (
            <Select.Option value={String(field.id)}>{field.name}</Select.Option>
          ))}
        </Select>
      </div>
    );
  } else if (fields[0].assessment_type_id === 4) {
    return (
      <div className={b("select")}>
        <label
          htmlFor={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          className={b("label")}
        >
          {fields[0].label}
        </label>
        <Select
          multiple={true}
          placeholder={fields[0].label}
          width={"max"}
          id={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          value={value}
          onUpdate={onUpdate}
        >
          {fields.map((field) => (
            <Select.Option value={String(field.id)}>{field.name}</Select.Option>
          ))}
        </Select>
      </div>
    );
  } else if (fields[0].assessment_type_id === 5) {
    return (
      <div className={b("text")}>
        <label
          htmlFor={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
          className={b("label")}
        >
          {fields[0].label}
        </label>
        <TextInput
          value={value[0]}
          onUpdate={(v) => {
            onUpdate([v]);
          }}
          placeholder={fields[0].label}
          id={`${uniqueId}_${String(fields[0].group_id)}_${String(fields[0].id)}`}
        />
      </div>
    );
  }
  return <></>;
};
