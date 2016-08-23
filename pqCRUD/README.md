# pqCRUD
The SQL code for the data structure:
```sql
CREATE TABLE userinfo
(
    uid serial NOT NULL,
    username character varying(100) NOT NULL,
    departname character varying(500) NOT NULL,
    Created date,
    CONSTRAINT userinfo_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);
```
If you have difficulties with the access right, here is an easy way to manage it using pgAdmin III:

1. Open **pgAdmin III**.
2. Double-click the **localhost server**, then type in your password.
3. Right-click **Login Roles**, select **New Login Role...**.
4. **Properties > Role name**: *penguin*.
5. **Definition > Password**: *penguin*, **Password (again)**: *penguin*.
6. **OK**.
7. Right-click **Databases**, select **New Database...**.
8. **Properties > Name**: *penguin*, **Owner**: *penguin*.
9. **OK**.
10. Right-click database **penguin**, select **Properties...**.
11. **Default Privileges > Tables > Role**: *public* (**ALL** is automatically checked), click **Add/Change**. Repeat for **Sequences**, **Functions**, and **Types**.
12. **OK**.
13. Click **SQL** button (tooltip: ==Execute Arbitrary SQL queries.==) in the toolbar.
14. Copy the SQL code above and paste it into **SQL Editor**.
15. Click **green triangle** button (tooltip: ==Execute query==).
16. Build and Run the Go code.
