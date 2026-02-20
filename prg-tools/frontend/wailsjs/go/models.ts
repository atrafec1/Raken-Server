export namespace domain {
	
	export class Job {
	    Name: string;
	    Number: string;
	
	    static createFrom(source: any = {}) {
	        return new Job(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Number = source["Number"];
	    }
	}
	export class Material {
	    BidNumber: string;
	    Name: string;
	    Unit: string;
	
	    static createFrom(source: any = {}) {
	        return new Material(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BidNumber = source["BidNumber"];
	        this.Name = source["Name"];
	        this.Unit = source["Unit"];
	    }
	}
	export class MaterialLog {
	    Job: Job;
	    Date: string;
	    Quantity: number;
	    Material: Material;
	
	    static createFrom(source: any = {}) {
	        return new MaterialLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Job = this.convertValues(source["Job"], Job);
	        this.Date = source["Date"];
	        this.Quantity = source["Quantity"];
	        this.Material = this.convertValues(source["Material"], Material);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class JobMaterialInfo {
	    Logs: MaterialLog[];
	    Materials: Material[];
	    Job: Job;
	    FromDate: string;
	    ToDate: string;
	
	    static createFrom(source: any = {}) {
	        return new JobMaterialInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Logs = this.convertValues(source["Logs"], MaterialLog);
	        this.Materials = this.convertValues(source["Materials"], Material);
	        this.Job = this.convertValues(source["Job"], Job);
	        this.FromDate = source["FromDate"];
	        this.ToDate = source["ToDate"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace dto {
	
	export class PayrollEntry {
	    EmployeeCode: string;
	    CurrentDate: string;
	    CraftLevel: string;
	    JobNumber: string;
	    Phase: string;
	    CostCode: string;
	    ChangeOrder: string;
	    RegularHours: number;
	    OvertimeHours: number;
	    PremiumHours: number;
	    Day: number;
	    EquipmentCode: string;
	    DownFlag: string;
	    SpecialPayType: string;
	    SpecialPayCode: string;
	    SpecialUnits: number;
	    SpecialRate: number;
	    CostCodeDivision: string;
	
	    static createFrom(source: any = {}) {
	        return new PayrollEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EmployeeCode = source["EmployeeCode"];
	        this.CurrentDate = source["CurrentDate"];
	        this.CraftLevel = source["CraftLevel"];
	        this.JobNumber = source["JobNumber"];
	        this.Phase = source["Phase"];
	        this.CostCode = source["CostCode"];
	        this.ChangeOrder = source["ChangeOrder"];
	        this.RegularHours = source["RegularHours"];
	        this.OvertimeHours = source["OvertimeHours"];
	        this.PremiumHours = source["PremiumHours"];
	        this.Day = source["Day"];
	        this.EquipmentCode = source["EquipmentCode"];
	        this.DownFlag = source["DownFlag"];
	        this.SpecialPayType = source["SpecialPayType"];
	        this.SpecialPayCode = source["SpecialPayCode"];
	        this.SpecialUnits = source["SpecialUnits"];
	        this.SpecialRate = source["SpecialRate"];
	        this.CostCodeDivision = source["CostCodeDivision"];
	    }
	}
	export class Warning {
	    Message: string;
	    WarningType: string;
	
	    static createFrom(source: any = {}) {
	        return new Warning(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Message = source["Message"];
	        this.WarningType = source["WarningType"];
	    }
	}
	export class PayrollEntryResult {
	    Entries: PayrollEntry[];
	    Warnings: Warning[];
	
	    static createFrom(source: any = {}) {
	        return new PayrollEntryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Entries = this.convertValues(source["Entries"], PayrollEntry);
	        this.Warnings = this.convertValues(source["Warnings"], Warning);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

