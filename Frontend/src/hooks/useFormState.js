import { useState } from 'react';

export default function useFormState(initialValues) {
  const [values, setValues] = useState(initialValues);
  const [errors, setErrors] = useState({});

  const handleChange = (event) => {
    const { name, value } = event.target;
    setValues((prev) => ({ ...prev, [name]: value }));
    setErrors((prev) => ({ ...prev, [name]: undefined }));
  };

  /** Marks every empty field in `requiredFields` as an error. Returns true when all are filled. */
  const validateRequired = (requiredFields) => {
    const missing = {};
    requiredFields.forEach((field) => {
      if (!String(values[field] ?? '').trim()) missing[field] = 'هذا الحقل مطلوب';
    });
    setErrors(missing);
    return Object.keys(missing).length === 0;
  };

  const reset = () => {
    setValues(initialValues);
    setErrors({});
  };

  return { values, errors, setErrors, handleChange, validateRequired, reset };
}
