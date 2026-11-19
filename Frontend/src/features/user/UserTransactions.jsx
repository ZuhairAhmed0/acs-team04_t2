import TransactionsView from '../../components/common/TransactionsView';
import { exportUserTransactions, fetchUserTransactions } from '../../services/user.service';

export default function UserTransactions() {
  return (
    <TransactionsView
      title="كل المعاملات"
      subtitle="سجل عمليات محفظتك، مع إمكانية التصفية وتنزيل ملف CSV."
      fetchTransactions={fetchUserTransactions}
      exportTransactions={exportUserTransactions}
    />
  );
}
