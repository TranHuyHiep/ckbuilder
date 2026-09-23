1. What is CKB
- Cell is the basic of CKB, similar to a cell in the human body.
- All the cells constitute the general state of the entire CKB blockchain.
- This process is the same as the Bitcoin UTXO.
- Unspent cells are *live* cells; spent cells are *dead* cells
- Unlike traditional UTXO, a cell can store any type of data. Each cell has a field called data, where you can put an unformatted string

2. What does a cell contain?
- Cells are obtained by the verification of the global consensus on the chain. Owning cells entails costs, since their storage space is limited.
- In Nervos CKB, 1 CKB is equal to 1 byte of storage space.

The entire cell data structure looks like this:
```bash
Cell: {
  capacity: HexString;
  lock: Script;
  type: Script;
  data: HexString;
}
```

The four fields are defined as follows:
- *capacity*: the space size of the cell, i.e. the integer number of native tokens represented by this cell, usually expressed in hexadecimal. The basic unit for capacity is shannon, 1 CKB equals 10**8 shannon.
- *lock*: a script, which is essentially the equivalent of a lock. We'll show you more details later.
- *type*: a script, same as the lock but for a different purpose.
- *data*: this field can store arbitrary bytes, which means any type of data

3. How to tell that you own a cell?
- If the cell is a box, the lock and type scripts are the two locks on the box.
- In essence, the scripts are a piece of code and parameters. When we try to consume a cell, the scripts run automatically and take the parameters and proof we submit (such as the signature) to determine if the locks of the cell can be unlocked. Once unlocked, it proves that we own and control the cell.

The script structure looks like this:
```bash
Script: {
  code_hash: HexString
  args: HexString
  hash_type: Uint8, there are three allowed values: {0: "data", 1: "type", 2: "data1"}
}
```

In these three fields, hash_type will be explained later, the other two are:
- *code_hash*: the hash of a certain piece of code
- *args*: the arguments that will be transferred to the code

3. Where is the code actually located?
- This dependency cell is called `CellDep` (dep cell).
- `code_hash` and `hash_type` fields are used to locate the code

4. What if the lock code is lost?
- Technically, the cell that contains the code of a lock should last as long as the chain, and no one can access this cell.

5. What is a transaction?
- Constructing a transaction is to destroy some cells and create some more
- transaction: inputs -> outputs

Transaction Rules
- The capacity summary of all the output cells must be less than the capacity summary of all the input cells:
- The difference in capacity between inputs and outputs, is the fee that the miner earns:

6. Difference between lock and type script
- The lock script is usually used to protect the ownership of the box, indicating who can unlock the box, while the type script is used to ensure that the cell follows certain application logic.
- The lock script is the gatekeeper, while the type script is the guardian.
- `Lock script`: In a transaction, the lock scripts run for all inputs by group.
- `Type script`: In a transaction, the type scripts run for all inputs and outputs by group.
*Note*: CKB does not run the script one by one. It first groups the inputs or outputs by script and runs the same script only once.


*Summary*
- CKB is essentially a chain of cells that are constantly being created and destroyed.
- A cell is a box that can be used to store all types of data.
- To own a cell and store data on-chain, you need tokens: `1 CKB = 1 Byte.`
- The byte size of the entire cell cannot exceed the value of the capacity field.
- To protect your cell, you must put a lock on the cell so that only you or someone you authorize can open it.
- A lock is essentially a piece of code that checks if the cell can be unlocked using arguments and user-provided signatures or proofs.
- The return value of 0 means that the lock was unlocked successfully, while any other value means the unlock attempt failed.
- The lock's code_hash and hash_type fields are used to locate code, which is stored in the data field of a dep cell.
- Each cell can carry two scripts, one is called lock script (default) and the other, type script (optional).
- In one transaction, the lock scripts run for all inputs by group, while the type scripts run for all inputs and outputs by group.

Summary 2
Congratulations! Now you are prepared for the practical tutorial!

Let's review all the concepts we have learned:

- CKB is essentially a chain of cells which are being created and destroyed over and over again.
- A cell is a box that can be used to store any type of data.
- To own a cell and store data on-chain, you need tokens: 1 CKB = 1 Byte.
- The byte size of the entire cell cannot exceed the value of the capacity field.
- To protect your cell, you must put a lock on the cell so that only you or someone you authorize can open it.
- A lock is essentially a piece of code that checks if the cell can be unlocked using arguments and user-provided signatures or proofs.
- The return value of 0 means that the lock was unlocked successfully, while any other value means the unlock attempt failed.
- The lock's code_hash and hash_type fields are used to locate code, which is stored in the data field of a dep cell.
- Each cell can carry two scripts, one is called lock script (default) and the other, type script (optional).
- In one transaction, the lock scripts run for all inputs by group, while the type scripts run for all inputs and outputs by group.
- The differences in the execution mechanism result in different uses for the two types of locks. Lock scripts are often used to protect the ownership of the cell. Type scripts often used to handle the cell transformation rules.
- Constructing a transaction is fundamentally about destroying some cells and creating some new ones.
- That's right, with the above theoretical knowledge, you're ready to hit the road.

![alt text](<Screenshot 2026-09-21 at 13.06.30.png>)